/**
 * @format
 */

import {AppRegistry, YellowBox} from 'react-native';
import App from './App';
import {name as appName} from './app.json';

YellowBox.ignoreWarnings([
    'Warning: componentWillReceiveProps is deprecated and will be removed in the next major version. Use static getDerivedStateFromProps instead.',
    'Warning: componentWillMount is deprecated and will be removed in the next major version. Use componentDidMount instead.',
]);
AppRegistry.registerComponent(appName, () => App);
