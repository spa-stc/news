<script lang="ts">
  import { pb, user } from "$lib/pocketbase";
  import { createMutation } from "@tanstack/svelte-query";
  import type { ClientResponseError } from "pocketbase";
  import type { Announcement } from "$lib/api/types";
  import toast from "svelte-french-toast";
  import { push } from "svelte-spa-router";
  import * as Card from "$lib/components/ui/card";
  import Button from "$lib/components/ui/button/button.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import Input from "$lib/components/ui/input/input.svelte";
  import Label from "$lib/components/ui/label/label.svelte";
  import Textarea from "$lib/components/ui/textarea/textarea.svelte";
  import * as Popover from "$lib/components/ui/popover";
  import {
    CalendarDate,
    DateFormatter,
    getLocalTimeZone,
    type DateValue,
    today,
  } from "@internationalized/date";
  import { type DateRange } from "bits-ui";
  import { CalendarIcon } from "lucide-svelte";
  import { RangeCalendar } from "$lib/components/ui/range-calendar";

  interface FormData {
    title: string;
    content: string;
    start_showing_at: string | undefined;
    finish_showing_at: string | undefined;
    user_id: string;
  }

  let userid = "";
  $: {
    if (!$user) {
      toast.error("You Must Be Logged In To Use This Feature");
      push("/login");
    } else {
      userid = $user.id;
    }
  }

  const mutation = createMutation<Announcement, ClientResponseError, FormData>({
    mutationFn: async (formdata): Promise<Announcement> => {
      const data = {
        title: formdata.title,
        content: formdata.content,
        author: formdata.user_id,
        approved: false,
        start_showing_at: formdata.start_showing_at,
        finish_showing_at: formdata.finish_showing_at,
      };

      return await pb.collection("announcements").create(data);
    },
    onSuccess: () => {
      toast.success("Submitted!");
    },
    onError: (err) => {
      toast.error(err.message);
    },
  });

  const df = new DateFormatter("en-US", {
    dateStyle: "full",
  });

  let now = today(getLocalTimeZone());

  let value: DateRange | undefined = {
    start: now,
    end: now.add({ days: 5 }),
  };

  let startValue: DateValue | undefined = undefined;

  let title = "";
  let content = "";
</script>

<div class="container-sm mx-auto mt-6 sm:mt-12">
  <Card.Root class="w-[350px] sm:w-[450px]">
    {#if $mutation.isSuccess}
      <Card.Header>
        <Card.Title>Thanks!</Card.Title>
        <Card.Description
          >Someone will review your submission, and it will appear on the
          student newsletter during the time you specified.</Card.Description
        >
      </Card.Header>
    {:else}
      <Card.Header>
        <Card.Title>Submit an announcement:</Card.Title>
        <Card.Description
          >Submit an announcement to appear on the student newsletter.</Card.Description
        >
      </Card.Header>
      <Card.Content class="space-y-2">
        <form on:submit|preventDefault>
          <FormItem>
            <Label>Title:</Label>
            <Input type="text" bind:value={title} />
          </FormItem>

          <FormItem>
            <Label>Content:</Label>
            <Textarea bind:value={content} />
          </FormItem>

          <FormItem>
            <Label>Display Range:</Label>
            <Popover.Root openFocus>
              <Popover.Trigger asChild let:builder>
                <Button
                  variant="outline"
                  class="w-full justify-start text-left font-normal {value
                    ? 'text-muted-foreground text-wrap'
                    : ''}"
                  builders={[builder]}
                >
                  <CalendarIcon class="mr-2 h-4 w-4" />
                  {#if value && value.start}
                    {#if value.end}
                      {df.format(value.start.toDate(getLocalTimeZone()))} - {df.format(
                        value.end.toDate(getLocalTimeZone())
                      )}
                    {:else}
                      {df.format(value.start.toDate(getLocalTimeZone()))}
                    {/if}
                  {:else if startValue}
                    {df.format(startValue.toDate(getLocalTimeZone()))}
                  {:else}
                    Pick a date
                  {/if}
                </Button>
              </Popover.Trigger>
              <Popover.Content class="w-auto p-0" align="start">
                <RangeCalendar
                  bind:value
                  bind:startValue
                  placeholder={value?.start}
                  initialFocus
                  numberOfMonths={1}
                />
              </Popover.Content>
            </Popover.Root>
          </FormItem>

          <Button
            class="mt-4"
            on:click={() => {
              $mutation.mutate({
                user_id: userid,
                title: title,
                content: content,
                start_showing_at: value?.start?.toString(),
                finish_showing_at: value?.start?.toString(),
              });
            }}>Submit</Button
          >
        </form>
      </Card.Content>
    {/if}
  </Card.Root>
</div>
