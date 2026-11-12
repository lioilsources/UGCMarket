package shop

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/lioilsources/RobloxMarket/api/internal/auth"
	"github.com/lioilsources/RobloxMarket/api/internal/db"
)

var ErrNotFound = errors.New("not found")

// IssueDownload issues a short-lived token for downloading a license's bundle. It verifies
// that the license belongs to the customer. The “My assets” page mints a token on
// every click, so downloads can be repeated without limit via the account.
func (s *Shop) IssueDownload(ctx context.Context, customerID, licenseID int64) (string, error) {
	lic, err := s.q.GetLicense(ctx, licenseID)
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && lic.CustomerID != customerID) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	token, hash := auth.NewToken()
	err = s.q.CreateDownloadToken(ctx, db.CreateDownloadTokenParams{TokenHash: hash, LicenseID: lic.ID,
		ExpiresAt: time.Now().Add(downloadTTL), MaxUses: downloadMaxUses})
	return token, err
}

type licenseData struct {
	Title, ProductURL, LicenseKey, Email, Issued string
	Files                                        []string
}

// ServeDownload streams a ZIP: the asset's files + LICENSE.txt + README.txt.
func (s *Shop) ServeDownload(w http.ResponseWriter, r *http.Request, token string) error {
	ctx := r.Context()
	licenseID, err := s.q.UseDownloadToken(ctx, auth.Hash(token))
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	lic, err := s.q.GetLicense(ctx, licenseID)
	if err != nil {
		return err
	}
	product, err := s.q.GetProductByID(ctx, lic.ProductID)
	if err != nil {
		return err
	}
	files, err := s.q.ListProductFiles(ctx, product.ID)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return ErrNotFound
	}
	email, err := s.customerEmail(ctx, lic.CustomerID)
	if err != nil {
		return err
	}

	names := make([]string, len(files))
	for i, f := range files {
		names[i] = fileName(product.Slug, f.Format)
	}
	data := licenseData{Title: product.Title, ProductURL: fmt.Sprintf("%s/p/%s", s.publicURL, product.Slug),
		LicenseKey: lic.LicenseKey, Email: email, Issued: lic.IssuedAt.UTC().Format("2006-01-02"), Files: names}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", product.Slug+".zip"))
	zw := zip.NewWriter(w)
	for i, f := range files {
		src, err := s.fs.Open(f.StorageKey)
		if err != nil {
			return err // headers already sent — the client gets a truncated zip
		}
		dst, err := zw.Create(product.Slug + "/" + names[i])
		if err == nil {
			_, err = io.Copy(dst, src)
		}
		src.Close()
		if err != nil {
			return err
		}
	}
	for name, tmpl := range map[string]string{"LICENSE.txt": lic.LicenseType + ".txt", "README.txt": "readme.txt"} {
		dst, err := zw.Create(product.Slug + "/" + name)
		if err != nil {
			return err
		}
		if _, err := dst.Write(render(tmpl, data)); err != nil {
			return err
		}
	}
	return zw.Close()
}

func fileName(slug, format string) string {
	if format == "png" {
		return slug + "_tex.png"
	}
	return slug + "." + format
}
