<script lang="ts">
    import Modal from "$lib/components/Modal.svelte";
    import {
        ListFilesInSnapshots,
        RestoreRepo,
        ChooseRepository,
        AddUpdateRepository,
        DeleteRepositorySetting,
        DeleteRepositoryFromDisk,
        GetRepoStats,
        CheckRepoErrors,
    } from "$lib/wailsjs/go/main/App";
    import { snapshots } from "$lib/store";
    import { myfetch_json, formatBytes } from "$lib/myfuncs";
    import FindResults from "$lib/components/FindResults.svelte";
    import SnapshotCrud from "$lib/components/SnapshotCrud.svelte";

    /* == Variables == */
    let modal_info = "";
    let modal_open = false;
    /* 
    0 : new record
    1 : update_id
    */
    let islem: number = 0;
    export let parent_snapshot_function: any = null;

    export let repo_add_modal = false;
    export let repo_form: any = {
        Name: null,
        Path: null,
        Password: null,
        Args: null,
    };
    /* == Variables End == */

    function delete_from_settings(): void {
        myfetch_json(DeleteRepositorySetting, $selected_repo_id).then(
            (r: any) => {
                repositories.set(r.names);
                repo_toggle();
            },
        );
    }
    /* 
    TODO: delete from disk
    */
    function delete_from_disk(): void {
        myfetch_json(DeleteRepositoryFromDisk, $selected_repo_id).then(
            (r: any) => {
                repositories.set(r.names);
                repo_toggle();
            },
        );
    }

    function choose_repository(): void {
        ChooseRepository().then((r: any) => {
            repo_form["Path"] = r;
        });
    }

    const repo_toggle = () => {
        repo_add_modal = !repo_add_modal;
    };

    function save_repo() {
        myfetch_json(
            AddUpdateRepository,
            $selected_repo_id,
            JSON.stringify(repo_form),
        ).then((r: any) => {
            repositories.set(r.names);
            repo_toggle();
            if (parent_snapshot_function) parent_snapshot_function();
        });
    }

    const get_stats = () => {
        myfetch_json(
            GetRepoStats,
            $selected_repo_id,
            JSON.stringify(repo_form),
        ).then((r: any) => {
            if (r) {
                toggleModal();
                modal_info =
                    "<strong>Snapshots Count</strong>: " +
                    r["snapshots_count"] +
                    "<br><strong>Total Files</strong>: " +
                    r["total_file_count"] +
                    "<br><strong>Disk Size</strong>: " +
                    formatBytes(r["disk_size"]) +
                    "<br><strong>Logical Size</strong>: " +
                    formatBytes(r["total_size"]);
            }
        });
    };

    const check_backup = () => {
        myfetch_json(
            CheckRepoErrors,
            $selected_repo_id,
            JSON.stringify(repo_form),
        ).then((r: any) => {
            if (r) {
                toggleModal();
                modal_info = String(r).replace(/(?:\r\n|\r|\n)/g, "<br>");
            }
        });
    };

    const toggleModal = () => {
        if (modal_open) {
            modal_info = "";
        }
        modal_open = !modal_open;
    };
</script>

{#if $selected_repo_id != -1}
    <FindResults />

    <div class="min-h-[200px] mt-4">
        {#if Array.isArray($snapshots)}
            <div class="space-y-2">
                {#each $snapshots as s_arr}
                    <details
                        class="group border border-gray-200 dark:border-gray-700 rounded-md bg-white dark:bg-gray-800"
                    >
                        <summary
                            class="flex justify-between items-center p-4 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700 list-none"
                        >
                            <div class="font-medium dark:text-white">
                                Snapshot Id: <span class="font-normal"
                                    >{s_arr["short_id"]}</span
                                >
                            </div>
                            <div class="text-gray-600 dark:text-gray-400">
                                At Time: {s_arr["time"]}
                            </div>
                        </summary>
                        <div
                            class="p-4 border-t border-gray-200 dark:border-gray-700"
                        >
                            <SnapshotCrud snapshot_info={s_arr} />
                        </div>
                    </details>
                {/each}
            </div>
            <hr class="my-4 dark:border-gray-700" />
        {/if}
    </div>

    <Modal isOpen={modal_open} toggle={toggleModal}>
        <div slot="header">Info</div>
        <div>{@html modal_info}</div>
    </Modal>

    <div class="flex flex-wrap gap-2 mt-4">
        <button class="btn btn-primary">Backup</button>
        <button class="btn btn-primary" on:click={check_backup}
            >Check Backup....</button
        >
        <button class="btn btn-primary">Show Difference</button>
        <button class="btn btn-primary">Forget</button>
        <button class="btn btn-primary" on:click={get_stats}>Stats</button>
        <button class="btn btn-primary">Unlock</button>
    </div>

    <Modal isOpen={repo_add_modal} toggle={repo_toggle} size="lg">
        <div slot="header">
            {#if islem == 0}
                Update / Delete Repository
            {/if}
        </div>

        <form class="space-y-4">
            <div>
                <label for="i4" class="form-label">Alias Name</label>
                <input
                    type="text"
                    class="form-control"
                    id="i4"
                    placeholder="like: Home Folder"
                    bind:value={repo_form["Name"]}
                />
            </div>

            <div>
                <label for="i1" class="form-label">Choose Repository</label>
                <div class="flex gap-2">
                    <input
                        type="text"
                        class="form-control"
                        id="i1"
                        bind:value={repo_form["Path"]}
                    />
                    <button
                        class="btn btn-primary whitespace-nowrap"
                        on:click|preventDefault={choose_repository}
                        >Choose Repository</button
                    >
                </div>
                <span class="text-sm text-gray-500 dark:text-gray-400"
                    >If you want to input like rclone you can write.</span
                >
            </div>

            <div>
                <label for="i2" class="form-label">Password</label>
                <input
                    type="password"
                    class="form-control"
                    id="i2"
                    bind:value={repo_form["Password"]}
                />
            </div>

            <div>
                <label for="i3" class="form-label">Other Args</label>
                <textarea
                    class="form-control"
                    id="i3"
                    placeholder="argclone=&quot;--tpslimit=10&quot; etc"
                    bind:value={repo_form["Args"]}
                />
                <span class="text-sm text-gray-500 dark:text-gray-400"
                    >You may put inputs as new line.</span
                >
            </div>
        </form>

        <div slot="footer">
            <div class="flex flex-wrap gap-2 w-full justify-between">
                <div class="flex gap-2">
                    <button
                        class="btn btn-danger"
                        on:click={delete_from_settings}>Delete Setting</button
                    >
                    <button class="btn btn-danger" on:click={delete_from_disk}
                        >Delete Repository From Disk ( Dangerous )</button
                    >
                </div>
                <div>
                    {#if $selected_repo_id == -1}
                        <button class="btn btn-success" on:click={save_repo}
                            >Save Repository Setting</button
                        >
                    {:else}
                        <button class="btn btn-success" on:click={save_repo}
                            >Update Repository Setting</button
                        >
                    {/if}
                </div>
            </div>
        </div>
    </Modal>
{/if}
