import { json, LoaderFunctionArgs } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";
import dayjs from "dayjs";
import Card from "~/components/card";
import { DayItem, DayItemOptional } from "~/components/dayitem";
import db from "~/db";
import { DaysDate } from "~/models/public/Days";

export default function Index() {
	const data = useLoaderData<typeof loader>();

	return (
		<>
			<h1 className="text-3xl font-bold my-4">Week Overview:</h1>
			<div>
				{data.days.map(day => (
					<Card key={day.date} extraClasses="my-2">
						<h1 className="text-lg font-bold mb-1">{dayjs(day.date).format("dddd, MMM D")}:</h1>
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
					</Card >
				))
				}
			</div >
		</>
	);
}

export async function loader({ }: LoaderFunctionArgs) {
	const today = dayjs()
	const start: Date = new Date(today.day(0).toString());
	const end: Date = new Date(today.day(6).toString());

	const days = await db.selectFrom('days')
		.selectAll().where('date', '>', start as DaysDate)
		.where('date', '<', end as DaysDate)
		.execute();

	return json(
		{
			days: days
		}
	);
};
