import {createStackNavigator} from 'react-navigation';
import MessagesScreen from './home/MessagesScreen';
import ThreadScreen from './home/ThreadScreen';

const HomeScreen = createStackNavigator({
    Messages: {screen: MessagesScreen},
    Thread: {screen: ThreadScreen},
});

export default HomeScreen;
