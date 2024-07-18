<script lang="ts">
    // J'ai juste choisi la facilité pour le css, faut que je le refasse histoire de pouvoir enlever les 3000 border à la place de mettre sur la div
    import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();
    export let page = 1;
    export let totalPages = 10;
    export let size = "md";
    export let textColor = "";
    export let borderColor = "border-slate-500";
    export let bgColor = "bg-slate-800";

    type Size = {
        [key: string]: string
    }

    const sizez: Size = {
        "sm": "text-sm px-2 py-1",
        "md": "px-4 py-2",
        "lg": "text-lg px-6 py-3"
    }

    const nextSize: Size = {
        "sm": "text-sm px-4 py-1",
        "md": "px-6 py-2",
        "lg": "text-lg px-8 py-3"
    }

    function previous() {
        dispatch('previous');
    }

    function next() {
        dispatch('next');
    }

    $: size = sizez[size];
    let size2 = nextSize[size];

</script>

<div class="flex justify-center rounded-md {textColor} mt-5">
    <div class="shadow-md flex">
        <button class="bg-panel {size} rounded-l-md border-[1px] {borderColor} {bgColor} transition duration-500 hover:bg-slate-700 active:bg-slate-600 active:transition-none" on:click={previous}>Previous</button>
        <div class="{size} border-y-[1px] {borderColor} {bgColor}">
            <span class="h-full">Page {page} of {totalPages}</span>
        </div>
        <button class="bg-panel {size2} rounded-r-md border-[1px] {borderColor} {bgColor} transition duration-500 hover:bg-slate-700 active:bg-slate-600 active:transition-none" on:click={next}>Next</button>
    </div>
</div>