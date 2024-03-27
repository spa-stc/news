import PocketBase from "pocketbase";
import { writable } from "svelte/store";

export const pb = new PocketBase("http://localhost:8090");

export const user = writable(pb.authStore.model);

pb.authStore.onChange((a) => {
  user.set(pb.authStore.model);
});
