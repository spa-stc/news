<script>
  import { Newspaper } from "lucide-svelte";
  import { link } from "svelte-spa-router";
</script>

<div class="flex gap-6 md:gap-10">
  <a use:link href="/" class="flex items-center space-x-2">
    <Newspaper />
    <span class="font-bold hidden sm:inline-block">SPA Student Newsletter</span>
  </a>
</div>
