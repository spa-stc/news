import Card from "./card";

export function DayItem(props: { title: string, content: string }) {
	return (
		<Card extraClasses="my-2">
			<h2 className="font-bold mb-1">{props.title}:</h2>
			<p>{props.content}</p>
		</Card>
	)
}

export function DayItemOptional(props: { title: string, content: string }) {
	if (props.content === "") {
		return (
			<DayItem {...props} />
		);
	} else { }
}
