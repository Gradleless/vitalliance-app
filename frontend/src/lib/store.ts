import { writable } from 'svelte/store';
import { main } from "./wailsjs/go/models";

export type PointageTime = {
    startPointageDate: string;
    date: string;
    hours: number;
    minutes: number;
    clientName: string;
    infos: string;
    endPointageDate: string;
};

export enum Month {
        January = "Janvier",
        February = "Février",
        March = "Mars",
        April = "Avril"  ,
        May = "Mai",
        June = "Juin",
        July = "Juillet",
        August = "Août",
        September = "Septembre",
        October = "Octobre",
        November = "Novembre",
        December = "Décembre"
}

const initialState = {
    loginUser: main.LoginResponse.createFrom(),
    pointages: new Map<Month, PointageTime[]>()
};

const store = writable(initialState);

export default store;