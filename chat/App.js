import React from 'react';
import {createAppContainer, createStackNavigator, createSwitchNavigator} from 'react-navigation';
import {ActivityIndicator, StatusBar, View} from 'react-native';
import {fromBottom, fromRight} from 'react-navigation-transitions';
import {DataService} from './DataService';
import Login from './auth/screens/Login';
import Signup from './auth/screens/Signup';
import Messages from './chat/screens/message/Messages';
import Create from './chat/screens/message/Create';
import Thread from './chat/screens/message/Thread';
import Settings from './chat/screens/settings/Settings';
import Edit from './chat/screens/settings/Edit';

//Buffer screen that waits to see if user logs in while restricting access to the rest of the app until status is determined
export class AuthLoadingScreen extends React.Component {
    constructor(props) {
        super(props);
        this._bootstrapAsync();
    }

    // Fetch the token from storage then navigate to our appropriate place
    _bootstrapAsync = async () => {
        // This will switch to the App screen or Auth screen and this loading
        // screen will be unmounted and thrown away.
        this.props.navigation.navigate((await DataService.isLoggedIn()) ? 'App' : 'Auth');
    };

    // Render any loading content that you like here
    render() {
        return (
            <View>
                <ActivityIndicator />
                <StatusBar barStyle="default" />
            </View>
        );
    }
}

//handles transitions for the create message page
const handleCustomTransition = ({scenes}) => {
    const prevScene = scenes[scenes.length - 2];
    const nextScene = scenes[scenes.length - 1];

    // Custom transitions go there
    if (prevScene && prevScene.route.routeName === 'Messages' && nextScene.route.routeName === 'Create') {
        return fromBottom(500);
    }
    return fromRight();
};

//Handles Messages screen, New message screen, and individual thread screen
const MessageStack = createStackNavigator(
    {
        Messages: {screen: Messages},
        New: {screen: Create},
        Thread: {screen: Thread},
    },
    {
        transitionConfig: nav => handleCustomTransition(nav),
    },
);

//Handles Settings and Edit screens
const SettingsStack = createStackNavigator({
    Settings: {screen: Settings},
    Edit: {screen: Edit},
});

//Switches between Messages and Settings
export const ChatNavigator = createStackNavigator({
    Messages: {
        screen: MessageStack,
        navigationOptions: {
            header: null,
        },
    },
    Settings: {
        screen: SettingsStack,
        navigationOptions: {
            header: null,
        },
    },
});

//Handles Login and Signup screens
export const AuthNavigator = createStackNavigator(
    {
        Login: {screen: Login},
        Signup: {screen: Signup},
    },
    {
        initialRouteName: 'Login',
    },
);

//Root app container, includes Chat and Auth Navigators
export default createAppContainer(
    createSwitchNavigator(
        {
            AuthLoading: AuthLoadingScreen,
            App: ChatNavigator,
            Auth: AuthNavigator,
        },
        {
            initialRouteName: 'AuthLoading',
        },
    ),
);
