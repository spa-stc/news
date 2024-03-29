<script lang="ts">
  import * as Avatar from "$lib/components/ui/avatar";
  import toast from "svelte-french-toast";
  import Button from "../ui/button/button.svelte";

  export let username: string;
</script>

<!-- <Button variant="outline" size="icon"> -->
<Avatar.Root>
  <Avatar.Image src="https://api.dicebear.com/8.x/thumbs/svg?seed={username}"
  ></Avatar.Image>
  <Avatar.Fallback>{username}</Avatar.Fallback>
</Avatar.Root>
<!-- </Button> -->
