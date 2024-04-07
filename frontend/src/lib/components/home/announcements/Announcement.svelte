<script lang="ts">
  import type { Announcement } from "$lib/api/types";
  import * as Card from "$lib/components/ui/card";
  import DOMPurify from "dompurify";

  export let data: Announcement;

  $: clean = DOMPurify.sanitize(data.content, { USE_PROFILES: { html: true } });
</script>

<Card.Root>
  <Card.Header>
    <Card.Title>{data.title}</Card.Title>
    <Card.Description>From: {data.expand.author.name}</Card.Description>
  </Card.Header>
  <Card.Content class="prose prose-slate prose-sm dark:prose-invert">
    {@html clean}
  </Card.Content>
</Card.Root>
