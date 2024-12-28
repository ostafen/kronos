import {FaRegTrashCan} from "react-icons/fa6";
import {useQueryClient} from "@tanstack/react-query";
import useDeleteSchedule from "@/hooks/use-delete-schedule.ts";
import DialogActionTrigger from "@/components/molecules/DialogActionTrigger/DialogActionTrigger.tsx";
import {useNavigate} from "react-router";

export default function DeleteScheduleTrigger(props: { id: string }) {
    const deleteSchedule = useDeleteSchedule();
    const queryClient = useQueryClient();
    const navigate = useNavigate();

    const handleDeleteSchedule = async (id: string) => {
        await deleteSchedule.mutateAsync(id);
        await queryClient.invalidateQueries({queryKey: ["schedules"]});
    };

    return (
        <DialogActionTrigger
            title="Delete schedule"
            colorPalette="red"
            variant="plain"
            p={0}
            aria-label="Delete schedule"
            dialogData={{
                title: "Delete schedule",
                content: <p>Do you really want to delete this schedule?</p>
            }}
            onSuccess={() => navigate('/')}
            onConfirm={() => handleDeleteSchedule(props.id)}
        >
            <FaRegTrashCan/>
            Delete schedule
        </DialogActionTrigger>
    );
}