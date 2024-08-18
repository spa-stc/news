import {
	Links,
	Meta,
	Outlet,
	Scripts,
	ScrollRestoration,
} from "@remix-run/react";
import "./tailwind.css";
import { MetaFunction } from "@remix-run/node";
import Navbar from "./components/navbar";

export const meta: MetaFunction = () => {
	return [
		{ title: "SPA SSN" },
	];
};

export function Layout({ children }: { children: React.ReactNode }) {
	return (
		<html lang="en">
			<head>
				<meta charSet="utf-8" />
				<meta name="viewport" content="width=device-width, initial-scale=1" />
				<Meta />
				<Links />
			</head>
			<body className="bg-gray-200">
				<Navbar />
				<div className="max-w-4xl mx-auto">
					{children}
				</div>
				<ScrollRestoration />
				<Scripts />
			</body>
		</html>
	);
}

export default function App() {
	return <Outlet />;
}
