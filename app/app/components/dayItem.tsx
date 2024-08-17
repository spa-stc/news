export function DayItem(props: { title: string, content: string }) {
	return (
		<div>
			<p>{props.title}</p>
			<p>{props.content}</p>
		</div>
	)
}

export function DayItemOptional(props: { title: string, content: string }) {
	if (props.content === "") {
		return (
			<DayItem {...props} />
		);
	} else { }
}
