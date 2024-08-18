import dayjs from "dayjs";

export function GetWeek(today: dayjs.Dayjs): string[] {
	const days: Array<string> = new Array();
	for (var i = 0; i < 7; i++) {
		days.push(today.day(i).format("YYYY-MM-DD"));
	}


	return days;
}
