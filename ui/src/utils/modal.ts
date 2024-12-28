import createSubject from "rx-subject";
import DialogData from "@/model/dialog-data.ts";


export const dialogClose$ = createSubject<void>();
export const dialogConfirm$ = createSubject<void>();
export const dialogReset$ = createSubject<void>();
export const dialogOpen$ = createSubject<DialogData>();