import {createContext, Dispatch, ReactNode, useReducer} from "react";

interface AppContextModel {
    state: AppState;
    dispatch: Dispatch<AppState>
}

const initialState: AppState = {
    animationStatus: "normal"
}

const AppContext = createContext<AppContextModel>({
    state: initialState,
    dispatch: () => initialState
});

interface AppState {
    animationStatus: "normal" | "played";
}

type AppReducer = (state: AppState, action: AppState) => AppState;

interface AppContextProviderProps extends AppState {
    children: (state: AppState) => ReactNode;
}

export const AppContextProvider = ({children, ...props}: AppContextProviderProps) => {
    const [state, dispatch] = useReducer<AppReducer>(
        (state, action) => ({...state, ...action}),
        props
    );

    return (
        <AppContext.Provider value={{state, dispatch}}>
            {children(state)}
        </AppContext.Provider>
    );
}

export default AppContext;
