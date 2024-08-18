export default function Card({ children, extraClasses }: { children: React.ReactNode, extraClasses?: string }) {
	return (
		<div className={`border-2 border-gray-300 rounded-md p-2 ${extraClasses}`}>
			{children}
		</div>
	)
}
