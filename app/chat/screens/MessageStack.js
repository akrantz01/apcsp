import {createStackNavigator} from 'react-navigation';
import MessagesScreen from './message/MessagesScreen';
import ThreadScreen from './message/ThreadScreen';

const MessageStack = createStackNavigator({
    Messages: {screen: MessagesScreen},
    Thread: {screen: ThreadScreen},
});

export default MessageStack;
