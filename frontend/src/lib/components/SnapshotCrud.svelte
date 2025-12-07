<script lang="ts">
  import Modal from "$lib/components/Modal.svelte";
  import {
    ListFilesInSnapshots,
    RestoreRepo,
    RestoreFilesInSnapshots,
  } from "$lib/wailsjs/go/main/App";
  import { selected_repo_id } from "$lib/store";
  import {
    myfetch_json,
    formatBytes,
    warning_notify,
    htmlEntities,
  } from "$lib/myfuncs";

  /* == Variables == */
  let modal_info = "";
  let modal_open = false;
  let modal_open_sm = false;
  let snapshot_files_list: any[] = [];
  let cb_form: any;
  export let snapshot_info: any = { id: null };

  const restore_snapshot = () => {
    myfetch_json(RestoreRepo, $selected_repo_id, snapshot_info.id).then(
      (r: any) => {
        if (r) {
          toggleModalSm();
          modal_info = String(r).replace(/(?:\r\n|\r|\n)/g, "<br>");
        }
      },
    );
  };

  const list_files_in_snapshot = () => {
    myfetch_json(
      ListFilesInSnapshots,
      $selected_repo_id,
      snapshot_info.id,
    ).then((r: any) => {
      if (r) {
        toggleModal();
        snapshot_files_list = r.slice(1); /* except first index */
      }
    });
  };

  const set_tag_snapshot = () => {
    myfetch_json(
      ListFilesInSnapshots,
      $selected_repo_id,
      snapshot_info.id,
    ).then((r: any) => {
      if (r) {
        console.log(r);
      }
    });
  };

  const toggleModal = () => {
    modal_open = !modal_open;
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
    let values = new Array();
    for (let i = 0; i < el.length; i++) {
      if (cb_form.elements[i].checked) values.push(cb_form.elements[i].value);
    }
    if (values.length > 0) {
      myfetch_json(
        RestoreFilesInSnapshots,
        $selected_repo_id,
        snapshot_info.id,
        JSON.stringify(values),
      ).then((r: any) => {
        toggleModalSm();
        modal_info = String(r).replace(/(?:\r\n|\r|\n)/g, "<br>");
        modal_info = htmlEntities(modal_info);
      });
    }
  };
</script>

<Modal isOpen={modal_open_sm} toggle={toggleModalSm}>
  <div>{@html modal_info}</div>
</Modal>

{#if snapshot_info.id != 0}
  <div class="flex flex-wrap gap-2 mb-2">
    <button class="btn btn-primary" on:click={restore_snapshot}
      >Restore Snapshot</button
    >
    <button class="btn btn-primary" on:click={list_files_in_snapshot}
      >List Files In Snapshot</button
    >
    <button class="btn btn-primary" on:click={set_tag_snapshot}
      >Tag Snapshot</button
    >
  </div>

  <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 text-sm dark:text-gray-300">
    {#each Object.keys(snapshot_info) as sp}
      <div><strong>{sp}</strong> : {snapshot_info[sp]}</div>
    {/each}
  </div>
{/if}

<Modal isOpen={modal_open} toggle={toggleModal} size="lg">
  {#if snapshot_files_list.length > 0}
    <button class="btn btn-primary mb-4" on:click={extract_files}
      >Extract Selected Files</button
    >
    <hr class="mb-4 dark:border-gray-700" />

    <form class="overflow-x-auto" bind:this={cb_form}>
      <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
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
              >Name</th
            >
            <th
              scope="col"
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
              >Path</th
            >
            <th
              scope="col"
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
              >Type</th
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
          {#each snapshot_files_list as tm, i}
            <tr class="hover:bg-gray-50 dark:hover:bg-gray-800">
              <td
                class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400"
                >{i + 1}</td
              >
              <td
                class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400"
                >{tm.name}</td
              >
              <td
                class="px-6 py-4 text-sm text-gray-500 dark:text-gray-400 break-all max-w-xs overflow-hidden"
                >{tm.path}</td
              >
              <td
                class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400"
                >{tm.type}</td
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
                  value={tm.path}
                  name={`cb` + i}
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </form>
  {/if}
</Modal>
