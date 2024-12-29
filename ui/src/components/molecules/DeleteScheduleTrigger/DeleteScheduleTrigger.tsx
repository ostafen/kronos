import {FaRegTrashCan} from "react-icons/fa6";
import {useQueryClient} from "@tanstack/react-query";
import useDeleteSchedule from "@/hooks/use-delete-schedule.ts";
import DialogActionTrigger from "@/components/molecules/DialogActionTrigger/DialogActionTrigger.tsx";
import {useNavigate} from "react-router";
import {PropsWithChildren} from "react";
import {IconButtonProps} from "@chakra-ui/react";

interface DeleteScheduleTriggerProps extends IconButtonProps {
    scheduleId: string | string[];
}

export default function DeleteScheduleTrigger(props: PropsWithChildren<DeleteScheduleTriggerProps>) {
    const {scheduleId, children, ...rest} = props;
    const deleteSchedule = useDeleteSchedule();
    const queryClient = useQueryClient();
    const navigate = useNavigate();

    const handleDeleteSchedule = async (id: string | string[]) => {
        const scheduleIds = typeof id === 'string' ? [id] : id;

        for (const id of scheduleIds) {
            await deleteSchedule.mutateAsync(id);
            await queryClient.invalidateQueries({queryKey: ["schedules"]});
        }
    };

    const content = scheduleId instanceof Array && scheduleId.length > 1 ?
        "Do you really want to delete these schedules?" :
        "Do you really want to delete this schedule?";

    return (
        <DialogActionTrigger
            title="Delete schedule"
            dialogData={{
                title: "Delete schedule",
                content: <p>{content}</p>
            }}
            onSuccess={() => navigate('/')}
            onConfirm={() => handleDeleteSchedule(scheduleId)}
            {...rest}
        >
            {children || (
                <>
                    <FaRegTrashCan/>
                    Delete schedule
                </>
            )}
        </DialogActionTrigger>
    );
}