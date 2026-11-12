-- Dev and test DB for the Mac (Docker doesn't run on the Mac). Production is the database "market".
-- The e2e test wipes the test DB, so its name must end in _test.
CREATE DATABASE market_dev OWNER market;
CREATE DATABASE market_test OWNER market;
