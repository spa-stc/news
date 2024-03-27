import Home from "./routes/Home.svelte";
import NotFound from "./routes/NotFound.svelte";
import { wrap } from "svelte-spa-router/wrap";

export const routes = {
  "/": Home,
  "*": NotFound,
};
