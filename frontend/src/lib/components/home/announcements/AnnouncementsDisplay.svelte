<script lang="ts">
  import { useAnnouncements } from "$lib/queries/announcements";
  import Announcement from "./Announcement.svelte";

  let page = 1;

  $: announcements = useAnnouncements(page, 30, "approved = true");
</script>

{#if $announcements.isPending}
  Pending
{:else if $announcements.isError}
  Error
{:else if $announcements.isSuccess}
  {#each $announcements.data.items as data}
    <Announcement {data} />
  {/each}
{/if}
