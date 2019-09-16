import React, {Component} from 'react';
import {GiftedChat} from 'react-native-gifted-chat';
import LinearGradient from 'react-native-linear-gradient';
import {StyleSheet} from 'react-native';

export default class ThreadScreen extends Component {
    static navigationOptions = ({navigation}) => {
        return {
            title: navigation.getParam('name', 'null'),
            header: null,
        };
    };

    state = {
        messages: [
            {
                _id: 1,
                text: 'Hello developer',
                createdAt: new Date(),
                user: {
                    _id: 2,
                    name: 'React Native',
                    avatar: 'https://placeimg.com/140/140/any',
                },
            },
        ],
    };

    onSend(messages = []) {
        this.setState(previousState => ({
            messages: GiftedChat.append(previousState.messages, messages),
        }));
    }

    render() {
        return (
            <LinearGradient style={styles.background} colors={['#00af3a', '#005baf']}>
                <GiftedChat
                    alignTop={true}
                    messages={this.state.messages}
                    onSend={messages => this.onSend(messages)}
                    user={{
                        _id: 1,
                    }}
                />
            </LinearGradient>
        );
    }
}

const styles = StyleSheet.create({
    background: {
        backgroundColor: '#EEEEEE',
        height: '100%',
        width: '100%',
    },
});
