import { Kysely, PostgresDialect } from "kysely";
import pkg from 'pg';
import Database from "./models/Database";
const { Pool } = pkg;

export default new Kysely<Database>(
	{
		dialect: new PostgresDialect({
			pool: new Pool({
				connectionString: process.env.DATABASE_URL,
			}),
		}),
	}
);
