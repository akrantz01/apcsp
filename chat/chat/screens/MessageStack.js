import {createStackNavigator} from 'react-navigation';
import MessagesScreen from './message/MessagesScreen';
import ThreadScreen from './message/ThreadScreen';
import NewScreen from './message/NewScreen';
import {fromRight, fromBottom} from 'react-navigation-transitions';

const handleCustomTransition = ({scenes}) => {
    const prevScene = scenes[scenes.length - 2];
    const nextScene = scenes[scenes.length - 1];

    // Custom transitions go there
    if (prevScene && prevScene.route.routeName === 'Messages' && nextScene.route.routeName === 'New') {
        return fromBottom(500);
    }
    return fromRight();
};

const MessageStack = createStackNavigator(
    {
        Messages: {screen: MessagesScreen},
        New: {screen: NewScreen},
        Thread: {screen: ThreadScreen},
    },
    {
        transitionConfig: nav => handleCustomTransition(nav),
    },
);

export default MessageStack;
