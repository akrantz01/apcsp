import {createStackNavigator} from 'react-navigation';
import {SettingsScreen} from './settings/SettingsScreen';
import {EditScreen} from './settings/EditScreen';

const SettingsStack = createStackNavigator({
    Settings: {screen: SettingsScreen},
    Edit: {screen: EditScreen},
});

export default SettingsStack;
