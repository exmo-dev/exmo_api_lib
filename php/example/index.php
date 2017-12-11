<?php

    require ('exmo.php');
    
    $exmo = new exmo (
        'K-0743a333e21fa0103ccc6025ed7bfe8b9529d822',
        'S-9c7b12ecb9eacd8a79309bc4c3e3886ad3621af1',
        'http://api.exmo.com/v1/'
    );
    
    print '<pre>';
    print_r ($exmo -> user_info ());
    print '</pre>';
