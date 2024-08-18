import { Link } from "@remix-run/react";
import { Menu } from "lucide-react";

export default function Navbar() {
	return (
		<div className="w-full px-4 py-4 border-b-2 border-gray-300 flex flex-row">
			<button className="border-2 border-gray-300 rounded-md">
				<Menu color="black" size={28} />
			</button>
			<div className="mx-auto"></div>
			<Link to="/" className="border-2 border-gray-300 rounded-md px-2 font-bold">Login</Link>
		</div >
	);
}
