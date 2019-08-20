import {createSwitchNavigator, createAppContainer} from 'react-navigation';
import App from './chat/ChatNavigator';
import {AuthNavigator} from './auth/AuthNavigator';
import {AuthLoadingScreen} from './auth/screens/AuthLoadingScreen';

// Implementation of MessageStack, OtherScreen, SignInScreen, AuthLoadingScreen
// goes here.

export default createAppContainer(
    createSwitchNavigator(
        {
            AuthLoading: AuthLoadingScreen,
            App: App,
            Auth: AuthNavigator,
        },
        {
            initialRouteName: 'AuthLoading',
        },
    ),
);
