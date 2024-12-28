import {FaRegTrashCan} from "react-icons/fa6";
import {IconButtonProps, useDisclosure} from "@chakra-ui/react";
import ConfirmDialog, {
    dialogClose$,
    dialogConfirm$,
    dialogReset$
} from "@/components/molecules/ConfirmDialog/ConfirmDialog.tsx";
import {take} from "rxjs/operators";
import {useQueryClient} from "@tanstack/react-query";
import useDeleteSchedule from "@/hooks/use-delete-schedule.ts";
import {PropsWithChildren} from "react";
import {Button} from "@/components/chakra/button.tsx";

interface DeleteScheduleTriggerProps extends IconButtonProps {
    scheduleId: string;
    onSuccess?: () => void;
}

export default function DeleteScheduleTrigger(props: PropsWithChildren<DeleteScheduleTriggerProps>) {
    const {scheduleId, children, onSuccess, ...iconButtonProps} = props;
    const deleteSchedule = useDeleteSchedule();
    const queryClient = useQueryClient();
    const {open, onOpen, onClose} = useDisclosure();

    const handleDeleteSchedule = async (id: string) => {
        await deleteSchedule.mutateAsync(id);
        await queryClient.invalidateQueries({queryKey: ["schedules"]});
    };

    const initDeleteFlow = async (scheduleId: string) => {
        onOpen();

        try {
            await new Promise((resolve, reject) => {
                dialogConfirm$
                    .pipe(take(1))
                    .subscribe(() => resolve({status: "confirmed"}));

                dialogClose$
                    .pipe(take(1))
                    .subscribe(() => reject({status: "canceled"}));
            });

            await handleDeleteSchedule(scheduleId);

            dialogReset$.next();

            onSuccess?.();
        } catch (error) {
            console.error(error);
        } finally {
            onClose();
        }
    };

    return (
        <>
            <ConfirmDialog
                title="Delete schedule"
                content={<p>Do you really want to delete this schedule?</p>}
                isOpen={open}
            />
            <Button
                title="Delete schedule"
                onClick={() => initDeleteFlow(scheduleId)}
                variant="ghost"
                p={0}
                aria-label="Delete schedule"
                {...iconButtonProps}
            >
                <FaRegTrashCan/>
                {children}
            </Button>
        </>
    )
}