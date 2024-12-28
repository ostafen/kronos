import {StrictMode} from 'react';
import {createRoot} from 'react-dom/client';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import {Provider} from './components/chakra/provider.tsx';
import {BrowserRouter, Route, Routes} from 'react-router';
import NewSchedulePage from './pages/NewSchedule.page.tsx';
import ScheduleListPage from './pages/ScheduleList.page.tsx';
import RootLayout from './components/layouts/Root.layout.tsx';
import ScheduleDetailPage from "@/pages/ScheduleDetail.page.tsx";
import ConfirmDialog from "@/components/molecules/ConfirmDialog/ConfirmDialog.tsx";

const queryClient = new QueryClient();

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <QueryClientProvider client={queryClient}>
            <Provider defaultTheme="light">
                <BrowserRouter basename="/web">
                    <Routes>
                        <Route element={<RootLayout/>}>
                            <Route index element={<ScheduleListPage/>}/>
                            <Route path="/new-schedule" element={<NewSchedulePage/>}/>
                            <Route path="/schedule/:scheduleId" element={<ScheduleDetailPage/>}/>
                        </Route>
                    </Routes>
                    <ConfirmDialog/>
                </BrowserRouter>
            </Provider>
        </QueryClientProvider>
    </StrictMode>,
);
