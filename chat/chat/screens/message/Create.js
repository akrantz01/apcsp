import React from 'react';
import {Component} from 'react';
import {View, StyleSheet} from 'react-native';
import LinearGradient from 'react-native-linear-gradient';
import {Card} from '../settings/Settings';

export default class Create extends Component {
    static navigationOptions = {
        title: 'Create Message',
        header: null,
    };

    render() {
        return (
            <LinearGradient colors={['#00af3a', '#005baf']}>
                <View style={styles.container}>
                    <Card
                        text={'Back'}
                        textColor={'#444444'}
                        iconName={'arrow-left'}
                        iconType={'feather'}
                        iconColor={'#444444'}
                        showArrow={false}
                        height={50}
                        onPress={() =>
                            this.props.navigation.dispatch({
                                type: 'Navigation/BACK',
                            })
                        }
                    />
                </View>
            </LinearGradient>
        );
    }
}

const styles = StyleSheet.create({
    container: {top: '5%', width: '100%', height: '100%'},
});
