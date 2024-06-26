<script lang="ts">
  import { useAnnouncements } from "$lib/queries/announcements";
  import Announcement from "./Announcement.svelte";
  import * as Pagination from "$lib/components/ui/pagination";
  import { ChevronRight, ChevronLeft } from "lucide-svelte";
  import toast from "svelte-french-toast";
  import ErrorResponse from "$lib/components/ErrorResponse.svelte";
  import { getDateQueryFields } from "./date";

  let page_num: number = 1;
  let perPage: number = 10;
  let siblingCount: number = 1;
  let totalItems = 0;

  $: announcements = useAnnouncements(page_num, perPage, {
    date: getDateQueryFields(),
    approved: true,
  });
  $: {
    if ($announcements.isSuccess) {
      totalItems = $announcements.data.totalItems;
    }

    if ($announcements.isError) {
      toast.error(`Something Went Wrong: ${$announcements.error}`);
    }
  }
</script>

{#if $announcements.isPending}
  <h1 class="text-xl mx-auto text-center my-12">Loading...</h1>
{:else if $announcements.isError}
  <ErrorResponse />
{:else if $announcements.isSuccess}
  {#each $announcements.data.items as data}
    <Announcement {data} />
  {/each}
{/if}
<Pagination.Root
  {perPage}
  {siblingCount}
  count={totalItems}
  let:pages
  let:currentPage
  class="mt-2"
  onPageChange={(num) => (page_num = num)}
>
  <Pagination.Content>
    <Pagination.Item>
      <Pagination.PrevButton>
        <ChevronLeft class="h-4 w-4" />
        <span class="hidden sm:block">Previous</span>
      </Pagination.PrevButton>
    </Pagination.Item>
    {#each pages as page (page.key)}
      {#if page.type === "ellipsis"}
        <Pagination.Item>
          <Pagination.Ellipsis />
        </Pagination.Item>
      {:else}
        <Pagination.Item>
          <Pagination.Link {page} isActive={currentPage === page.value}>
            {page.value}
          </Pagination.Link>
        </Pagination.Item>
      {/if}
    {/each}
    <Pagination.Item>
      <Pagination.NextButton>
        <span class="hidden sm:block">Next</span>
        <ChevronRight class="h-4 w-4" />
      </Pagination.NextButton>
    </Pagination.Item>
  </Pagination.Content>
</Pagination.Root>
