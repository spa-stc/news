const { makeKyselyHook } = require("kanel-kysely");

const outputPath = './app/models';

/** @type {import('../src/Config').default} */
module.exports = {
	connection: {
		connectionString: process.env.DATABASE_URL,
		outputPath,
		resolveViews: true,
		preDeleteOutputFolder: true,
	},
	outputPath,

	preRenderHooks: [makeKyselyHook()],
}
