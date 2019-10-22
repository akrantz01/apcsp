/**
 * @format
 */

import {AppRegistry, YellowBox} from 'react-native';
import App from './App';
import {name as appName} from './app.json';

YellowBox.ignoreWarnings([
    'Warning: componentWillReceiveProps is deprecated and will be removed in the next major version. Use static getDerivedStateFromProps instead.',
    'Warning: componentWillMount is deprecated and will be removed in the next major version. Use componentDidMount instead.',
    'Remote debugger is in a background tab which may cause apps to perform slowly. Fix this by foregrounding the tab (or opening it in a separate window).',
]);
AppRegistry.registerComponent(appName, () => App);
