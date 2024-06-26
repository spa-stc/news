import Home from "./routes/Home.svelte";
import NotFound from "./routes/NotFound.svelte";
import { wrap } from "svelte-spa-router/wrap";

export const routes = {
  "/": Home,
  "/login": wrap({
    asyncComponent: () => import("./routes/LoginSignup.svelte"),
  }),
  "/post-signup": wrap({
    asyncComponent: () => import("./routes/PostSignup.svelte"),
  }),
  "/submit": wrap({
    asyncComponent: () => import("./routes/SubmitAnnouncement.svelte"),
  }),
  "*": NotFound,
};
