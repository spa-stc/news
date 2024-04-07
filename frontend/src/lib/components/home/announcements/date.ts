import { today, getLocalTimeZone } from "@internationalized/date";

export function getDateQueryFields(): string {
  let date = today(getLocalTimeZone());

  return date.toString();
}
