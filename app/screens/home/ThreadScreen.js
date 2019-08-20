import {Component} from 'react';

export default class ThreadScreen extends Component {
    static navigationOptions = ({navigation}) => {
        return {
            title: navigation.getParam('name', 'null'),
        };
    };

    render() {
        return null;
    }
}
