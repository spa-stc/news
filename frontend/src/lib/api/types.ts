import PocketBase, { RecordService } from "pocketbase";

interface User {
  id: string;
  username: string;
  email: string;
  name: string;
  role: "admin" | "student";
}

interface Announcement {
  id: string;
  title: string;
  content: string;
  author: string;
  approved: boolean;
}

// Pocket base instance, with API types.
export interface TypedPocketBase extends PocketBase {
  collection(idOrName: "users"): RecordService<User>;
  collection(idOrName: "announcements"): RecordService<Announcement>;
}
