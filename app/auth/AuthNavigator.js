import {createStackNavigator} from 'react-navigation';
import {Login} from './screens/Login';
import {Signup} from './screens/Signup';

export const AuthNavigator = createStackNavigator({
    Login: {screen: Login},
    Signup: {screen: Signup},
});
