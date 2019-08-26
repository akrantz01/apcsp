import {createStackNavigator} from 'react-navigation';
import {SettingsScreen} from './settings/SettingsScreen';

const SettingsStack = createStackNavigator({
    Settings: {screen: SettingsScreen},
});

export default SettingsStack;
