import type { Announcement } from "$lib/api/types";
import { pb } from "$lib/pocketbase";
import { createQuery } from "@tanstack/svelte-query";
import type { ClientResponseError, ListResult } from "pocketbase";

const fetchAnnouncements = async (
  page: number,
  query: string,
  per_page: number
): Promise<ListResult<Announcement>> => {
  const result = await pb.collection("announcements").getList(page, per_page, {
    filter: query,
    expand: "author",
  });

  return result;
};

export const useAnnouncements = (
  page: number,
  per_page: number,
  query: string
) =>
  createQuery<ListResult<Announcement>, ClientResponseError>({
    queryKey: ["announcements", page, per_page, query],
    queryFn: () => fetchAnnouncements(page, query, per_page),
  });
