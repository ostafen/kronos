import createSubject from 'rx-subject';
import DialogData from '@/model/dialog-data.ts';

export const closeDialog$ = createSubject<void>();
export const dialogConfirm$ = createSubject<void>();
export const dialogReset$ = createSubject<void>();
export const dialogOpen$ = createSubject<DialogData>();
export const isDialogOpen$ = createSubject<boolean>();
