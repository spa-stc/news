import { json, LoaderFunctionArgs } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";
import dayjs from "dayjs";
import { DayItem, DayItemOptional } from "~/components/dayItem";
import db from "~/db";
import { GetWeek } from "~/utils/time";

export default function Index() {
	const data = useLoaderData<typeof loader>();

	return (
		<div>
			{data.days.map(day => (
				<div key={day.date}>
					<h1>{dayjs(day.date).format("dddd, MMM D")}</h1>
					<DayItem title="Lunch" content={day.lunch} />
					<DayItemOptional title="Notes" content={day.notes} />
					<DayItemOptional title="Rotation Day" content={day.rotation_day} />
					<DayItemOptional title="X Period" content={day.x_period} />
					<DayItemOptional title="Location" content={day.location} />
					<DayItemOptional title="Grade 9" content={day.grade_9} />
					<DayItemOptional title="Grade 10" content={day.grade_10} />
					<DayItemOptional title="Grade 11" content={day.grade_11} />
					<DayItemOptional title="Grade 12" content={day.grade_12} />
					<DayItemOptional title="AP Info" content={day.ap_info} />
					<DayItemOptional title="CC Info" content={day.cc_info} />
				</div>
			))}
		</div>
	);
}

export async function loader({ }: LoaderFunctionArgs) {
	const dates = GetWeek(dayjs());

	const days = await db.selectFrom('days').selectAll().where("date", "in", dates).execute();

	return json(
		{
			days: days
		}
	);
};
