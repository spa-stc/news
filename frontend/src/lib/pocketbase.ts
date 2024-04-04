import PocketBase from "pocketbase";
import { writable } from "svelte/store";
import type { TypedPocketBase } from "./api/types";

export const pb = new PocketBase(
  import.meta.env.VITE_PB_URL
) as TypedPocketBase;

export const user = writable(pb.authStore.model);

pb.authStore.onChange((a) => {
  user.set(pb.authStore.model);
});
