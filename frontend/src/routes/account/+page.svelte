<script lang="ts">
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import store from "$lib/store";
    import { GetPointage, CalculateTimeHorodatage } from "../../lib/wailsjs/go/main/App";
    import type { PointageTime, Month } from "$lib/store";
    import { get } from 'svelte/store';
    import { slide } from 'svelte/transition';
    import Pagination from '$lib/components/Pagination.svelte';
    
    onMount(() => {
        const value = get(store);
        if (!value.loginUser) {
            goto("/");
        } else {
            GetPointage(value.loginUser.data["token"]).then(async (response) => {
            let pointages = await CalculateTimeHorodatage(response);
            const convertedPointages = new Map<Month, PointageTime[]>();
            Object.entries(pointages)
                .sort(([key1], [key2]) => {
                    const [month1, year1] = key1.split(" ");
                    const [month2, year2] = key2.split(" ");
                    // j'aurais dû utiliser un objet Date pour comparer les dates, avant de mettre en français
                    const monthOrder = ['janvier', 'février', 'mars', 'avril', 'mai', 'juin', 'juillet', 'août', 'septembre', 'octobre', 'novembre', 'décembre'];
                    if (year1 === year2) {
                        return monthOrder.indexOf(month1.toLowerCase()) - monthOrder.indexOf(month2.toLowerCase());
                    } else {
                        return year1.localeCompare(year2);
                    }
                })
                .forEach(([key, value]) => {
                    convertedPointages.set(key as Month, value as PointageTime[]);
                });
            store.set({ loginUser: value.loginUser, pointages: convertedPointages });
            let convertedPointagesKeys = Array.from(convertedPointages.keys());
            selectedMonth = convertedPointagesKeys[convertedPointagesKeys.length - 1];
            });
        }
    });

    let pointages: Map<Month, PointageTime[]> = new Map<Month, PointageTime[]>();
    store.subscribe((value) => {
        pointages = value.pointages;
    });
    
    let selectedMonth: Month;

    function calculateTotalTime(month: Month): string {
        let hours = 0;
        let minutes = 0;
        pointages.get(month)?.forEach((pointage) => {
            hours += pointage.hours;
            minutes += pointage.minutes;
            if(minutes >= 60) {
                hours += 1;
                minutes -= 60;
            }
        });
        return `${hours}h${minutes}`;
    }

    let totalTime: string = "0h0";

    $: {
        if(pointages.get(selectedMonth)) {
            totalTime = calculateTotalTime(selectedMonth);
        }
    }

    function calculatePercentage(): number {
        const previousMonthIndex = Array.from(pointages.keys()).indexOf(selectedMonth) - 1;
        if (previousMonthIndex < 0) {
            return 0;
        }
        const previousMonth = Array.from(pointages.keys())[previousMonthIndex];
        const previousMonthHours = parseInt(calculateTotalTime(previousMonth).split("h")[0]);
        const currentMonthHours = parseInt(calculateTotalTime(selectedMonth).split("h")[0]);
        const percentage = Math.round(((currentMonthHours - previousMonthHours) / previousMonthHours) * 100);
        return percentage;
    }

    let percentage: number = 0;
    let percentageSign: string = "";

    $: {
        if (pointages.get(selectedMonth)) {
            percentage = calculatePercentage();
            percentageSign = percentage >= 0 ? "+" : "";
        }
    }

    let firstName: string;
    let lastName: string;
    $: {
        const value = get(store).loginUser;
        if (value.success) {
            firstName = value.data["firstName"];
            lastName = value.data["lastName"];
        }
    }

    $: page = 1;
    $: totalPages = Math.ceil((pointages.get(selectedMonth)?.length ?? 0) / 10);
    
    function nextPage() {
        if (page < totalPages) {
            page += 1;
        }
    }

    function previousPage() {
        if (page >= 1) {
            page -= 1;
        }
    }
</script>

<div class="flex p-5 justify-center min-h-screen">
    <div class="w-full max-w-4xl p-6 bg-white rounded-lg">
        <div class="space-y-2 text-center">
            <h2 class="font-semibold text-2xl">Bienvenue, {firstName} {lastName} ! Pointages</h2>
            <p class="text-gray-500 dark:text-gray-400 text-sm">Liste de vos pointages</p>
        </div>

        <div class="grid gap-6 mt-10 md:mt-12 lg:mt-14">
            <div class="rounded-lg border bg-card text-card-foreground shadow-sm">
                <div class="flex flex-col space-y-1.5 p-6">
                    <h3 class="whitespace-nowrap text-2xl font-semibold leading-none tracking-tight">Vos stats</h3>
                </div>
                <div class="p-6">
                    <div class="grid md:grid-cols-3 gap-6">
                        <div>
                            <div class="text-4xl font-bold">{totalTime}</div>
                            <div class="text-muted-foreground mr-5">d'heures réalisées en {selectedMonth}</div>
                        </div>
                        <div>
                            <div class="text-4xl font-bold">{pointages.get(selectedMonth)?.length ?? 0}</div>
                            <div class="text-muted-foreground">Pointages totaux</div>
                        </div>
                        <div>
                            <div class="text-4xl font-bold">{percentageSign}{percentage}%</div>
                            <div class="text-muted-foreground">Heures comparé au mois précédent</div>
                        </div>
                    </div>
                </div>
            </div>
        <div class="mt-4">
            <select class="mb-4" bind:value={selectedMonth}>
            {#each Array.from(pointages.keys()) as month}
                <option value={month}>{month}</option>
            {/each}
            </select>
            {#if pointages.get(selectedMonth)}
            {#each pointages.get(selectedMonth)?.slice((page - 1) * 10, page * 10) ?? [] as pointage, index (pointage)}
                <div class="flex items center justify-between border-b border-gray-200 py-2" in:slide|global={{ delay: 350 }} out:slide|global>
                <div class="flex items-center">
                    <div class="flex flex-col">
                    <span class="text-sm text-gray-500 dark:text-gray-400">{new Date(pointage.date).toLocaleDateString()}</span>
                    <span class="text-sm text-gray-500 dark:text-gray-400">{pointage.hours > 0 ? `${pointage.hours}h` : ""}{pointage.minutes > 0 ? `${pointage.minutes}m` : ""}{pointage.seconds > 0 ? `${pointage.seconds}s` : ""}</span>
                    </div>
                </div>
                <div class="flex items-center">
                    <span class="text-sm text-gray-500 dark:text-gray-400">{pointage.clientName}</span>
                </div>
                </div>
            {/each}
            {#if pointages.get(selectedMonth)?.length ?? 0 > 10}
                <Pagination {page} {totalPages} size="md" on:next={nextPage} on:previous={previousPage} bgColor="bg-black" textColor="text-white" />
            {/if}
            {:else}
            <div class="text-center text-gray-500 dark:text-gray-400">Aucun pointage</div>
            {/if}
        </div>
        </div>
    </div>
</div>