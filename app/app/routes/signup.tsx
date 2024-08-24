import Card from "~/components/card";
import { z, ZodError } from 'zod';
import { Form, isRouteErrorResponse, json, useActionData, useRouteError } from "@remix-run/react";
import { Link } from "@remix-run/react";
import { ActionFunctionArgs } from "@remix-run/node";
import { zx } from "zodix";
import db from "~/db";
import * as argon2 from 'argon2';

const signupSchema = z.object({
	name: z.string().max(255, "Name must be less than 255 characters.").min(5, "Name must be greater than 5 characters."),
	email: z.string().email({ message: "Invalid Email Address" }),
	password: z.string().max(100, "Password must contain less than 100 characters.").min(6, "Password length must be greater than 6 characters."),
})

export default function Signup() {
	const data = useActionData<typeof action>();
	if (data?.success) {
		<Card>
			<h1 className="font-bold text-3xl">Success! You should be receiving a confirmation email shortly.</h1>
		</Card>
	}

	return (
		<div className="my-2">
			<Card>

				{data?.overallError && <div className="p-2 bg-red-100 border-2 border-gray-700 rounded-md">{data.overallError}</div>}
				<Form method="post" className="flex flex-col space-y-2">
					<h1 className="text-2xl font-bold">Signup</h1>
					<h2 className="font-bold">Name:</h2>
					<input className="border-2 border-gray-600 rounded-md pl-2" name="name" type="text" />
					{data?.nameError && <div className="p-2 bg-red-100 border-2 border-gray-700 rounded-md">{data.nameError}</div>}
					<h2 className="font-bold">Email:</h2>
					<input className="border-2 border-gray-600 rounded-md pl-2" name="email" type="text" />
					{data?.emailError && <div className="p-2 bg-red-100 border-2 border-gray-700 rounded-md">{data.emailError}</div>}
					<h2 className="font-bold">Password:</h2>
					<input className="border-2 border-gray-600 rounded-md pl-2" name="password" type="password" />
					{data?.passwordError && <div className="p-2 bg-red-100 border-2 border-gray-700 rounded-md">{data.passwordError}</div>}

					<div className="flex flex-row space-x-2">
						<button type="submit" className="rounded-md text-gray-50 bg-gray-700 px-2 font-bold py-1">Sign Up</button>
						<Link to="/login" className="self-baseline rounded-md border-2 border-gray-700 px-2 font-bold py-1">Login Instead</Link>
					</div>
				</Form>
			</Card>
		</div>
	)
}

export function ErrorBoundary() {
	const error = useRouteError();

	if (isRouteErrorResponse(error)) {
		return (
			<div>
				<h1 className="font-bold text-3xl">Oops:</h1>
				<p className="">Status: {error.status}</p>
				<p>{error.data.message}</p>
			</div>
		)
	}
}

function errorAtPath(error: ZodError, path: string) {
	return error.issues.find((issue) => issue.path[0] === path)?.message;
}

export async function action({ request }: ActionFunctionArgs) {
	const result = await zx.parseFormSafe(request, signupSchema);
	if (result.success) {
		const { name, email, password } = result.data;


		try {
			await db.insertInto("users").values(
				{
					name: name,
					password_hash: await argon2.hash(password),
					email: email,
				}
			).execute()
		} catch (error: any) {
			console.log(error);

			return json({
				success: false,
				overallError: "Something went wrong, please try again.",
				emailError: null,
				nameError: null,
				passwordError: null,
			})
		}

		return json({
			success: true,
			emailError: null,
			nameError: null,
			passwordError: null,
			overallError: null
		})
	}


	return json({
		success: false,
		emailError: errorAtPath(result.error, "email"),
		nameError: errorAtPath(result.error, "name"),
		passwordError: errorAtPath(result.error, "password"),
		overallError: null
	})
}
