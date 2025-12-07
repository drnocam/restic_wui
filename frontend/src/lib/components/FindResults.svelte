<script lang="ts">
  import {
    formatBytes,
    myfetch_json,
    warning_notify,
    htmlEntities,
  } from "$lib/myfuncs";
  import { selected_repo_id } from "$lib/store";
  import Modal from "$lib/components/Modal.svelte";
  import {
    SearchInRepo,
    RestoreFilesInSnapshots,
  } from "$lib/wailsjs/go/main/App";

  let modal_info = "";
  let modal_open_sm = false;
  let find_data: any[] = [];
  let search_value: any = null;
  let cb_form: any;
  let checked_items = [];

  const search_in = () => {
    myfetch_json(SearchInRepo, $selected_repo_id, search_value).then(
      (r: any) => {
        if (r) {
          find_data = r;
        }
      },
    );
  };

  const toggleModalSm = () => {
    if (modal_open_sm) {
      modal_info = "";
    }
    modal_open_sm = !modal_open_sm;
  };

  const extract_files = () => {
    if (!cb_form) {
      warning_notify("No Results Yet");
      return;
    }
    let el = cb_form.elements;
    if (el.length == 0) {
      warning_notify("Please Select Items First");
      return;
    }
    let values: any = {};
    for (let i = 0; i < el.length; i++) {
      if (cb_form.elements[i].checked) {
        if (values[cb_form.elements[i].dataset.snapshot] === undefined) {
          values[cb_form.elements[i].dataset.snapshot] = [];
        }
        values[cb_form.elements[i].dataset.snapshot].push(
          cb_form.elements[i].value,
        );
      }
    }

    if (Object.keys(values).length > 0) {
      Object.keys(values).forEach((key) => {
        myfetch_json(
          RestoreFilesInSnapshots,
          $selected_repo_id,
          key,
          JSON.stringify(values[key]),
        ).then((r: any) => {
          modal_open_sm = true;
          modal_info = String(r).replace(/(?:\r\n|\r|\n)/g, "<br>");
          modal_info = htmlEntities(modal_info);
        });
      });
    }
  };
</script>

<Modal isOpen={modal_open_sm} toggle={toggleModalSm}>
  <div>{@html modal_info}</div>
</Modal>

<div class="flex gap-2 mb-4">
  <input
    bind:value={search_value}
    type="text"
    class="form-control"
    placeholder="You may search by * symbols"
    aria-describedby="search-button"
    on:keyup={(e) => {
      if (e.key === "Enter") {
        search_in();
      }
    }}
  />
  <button
    class="btn btn-primary whitespace-nowrap"
    id="search-button"
    on:click={search_in}>Search</button
  >
</div>

<button class="btn btn-primary mb-4" on:click={extract_files}
  >Extract Selected Files</button
>
<hr class="mb-4 dark:border-gray-700" />

{#if find_data.length > 0}
  <form class="overflow-x-auto" bind:this={cb_form}>
    {#each find_data as t}
      <h6 class="font-bold mt-4 dark:text-white">Snapshot: {t.snapshot}</h6>
      <span class="text-sm text-gray-500 dark:text-gray-400">Hits {t.hits}</span
      >

      <table
        class="min-w-full divide-y divide-gray-200 dark:divide-gray-700 mt-2"
      >
        <thead class="bg-gray-50 dark:bg-gray-800">
          <tr>
            <th
              scope="col"
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
              >#</th
            >
            <th
              scope="col"
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
              >User</th
            >
            <th
              scope="col"
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
              >Path</th
            >
            <th
              scope="col"
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
              >Size</th
            >
            <th
              scope="col"
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
              >Created Time</th
            >
            <th
              scope="col"
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
              >Select</th
            >
          </tr>
        </thead>
        <tbody
          class="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700"
        >
          {#each t.matches as tm, i}
            <tr class="hover:bg-gray-50 dark:hover:bg-gray-800">
              <td
                class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400"
                >{i + 1}</td
              >
              <td
                class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400"
                >{tm.user}</td
              >
              <td
                class="px-6 py-4 text-sm text-gray-500 dark:text-gray-400 break-all"
                >{tm.path}</td
              >
              <td
                class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400"
                >{formatBytes(tm.size)}</td
              >
              <td
                class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400"
                >{new Date(Date.parse(tm.ctime)).toLocaleString("tr-TR")}</td
              >
              <td
                class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400"
              >
                <input
                  type="checkbox"
                  data-snapshot={t.snapshot}
                  value={tm.path}
                  name={`cb` + i}
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    {/each}
  </form>
{/if}
