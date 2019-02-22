
app = new Vue({
    el: '.app',
    data: {
        flag: ''
    },
    methods: {
        signout: function(e) {
            console.log('Signing out');
            $cookies.remove('auth');
            $cookies.remove('username');
            $cookies.remove('password');
        },
        hint: function(index) {
            console.log('Requesting hint');
            console.log(index);
            location.href = '/app/hint?index=' + index;
        },
        submit: function() {
            console.log('Submitting a flag');
            console.log(this.flag);
            location.href = '/app/submit?flag=' + this.flag;
        },
        reset: function() {
            console.log('Reseting all progress');
            location.href = '/app/reset';
        }
    },
});

