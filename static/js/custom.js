$(document).ready(function () {
    /* **** sticky **** */
    $(window).scroll(function () {
        if ($(this).scrollTop() > 150) {
            $("header").addClass("nav-new");
        } else {
            $("header").removeClass("nav-new");
        }
    });
    /* **** sticky **** */

    /* **** Slider ***** */
    $(".review-slider").slick({
        arrows: true,
        loop: true,
        dots: false,
        autoplay: true,
        autoplasySpeed: 1500,
        speed: 1000,
        infinite: true,
        slidesToShow: 2,
        slidesToScroll: 1,
        responsive: [
            {
                breakpoint: 1600,
                settings: {
                    slidesToShow: 2,
                },
            },
            {
                breakpoint: 1199,
                settings: {
                    slidesToShow: 1,
                },
            },
            {
                breakpoint: 991,
                settings: {
                    slidesToShow: 1,
                },
            },
            {
                breakpoint: 767,
                settings: {
                    slidesToShow: 1,
                },
            },
            {
                breakpoint: 575,
                settings: {
                    slidesToShow: 1,
                },
            },
            {
                breakpoint: 447,
                settings: {
                    slidesToShow: 1,
                },
            },
        ],
    });
    /* ***** End Slider **** */

     /* **** Slider ***** */
    $(".oblique-slider").slick({
        arrows: true,
        loop: true,
        dots: false,
        autoplay: false,
        autoplaySpeed: 1150,
        speed: 1000,
        infinite: false,
        slidesToShow: 1,
        slidesToScroll: 1,
    });
    /* ***** End Slider **** */

    /* **** scrollIt ***** */
    $(function () {
        $.scrollIt({
            upKey: 38,
            downKey: 40,
            easing: "linear",
            scrollTime: 600,
            activeClass: "active",
            onPageChange: null,
            topOffset: 0,
        });
    });
    /* **** End scrollIt ***** */

    /* **** Banner Slider **** */
    $(document).ready(function () {
        //save boolean
        var pause = false;
        //save items that with number
        var item=  $('.select-item');
        //save blocks
        var block=  $('.bg-block');
        //variable for counter
        var k =0;
          
          
        //interval function works only when pause is false
        setInterval(function () {
            if (!pause) {
                var $this = item.eq(k);
                  
                if (item.hasClass('active'))  {
                    item.removeClass('active');
                };
                  
                block.removeClass('active').eq(k).addClass('active');
                $this.addClass('active');
                //increase k every 1.5 sec
                k++;
                //if k more then number of blocks on page
                if (k >= block.length ) {
                    //rewrite variable to start over
                    k = 0;
                }
            }
            //every 1.5 sec
        }, 10000);false


        item.hover(function () {
            //remove active class from all except this
            $(this).siblings().removeClass("active");
            $(this).addClass('active');
            //remove active class from all
            block.removeClass('active');
                
            //add active class to block item which is accounted as index cliked item
            block.eq($(this).index()).addClass('active');
            //on hover stop interval function
            pause = true;
        }, function () {
            //when hover event ends, start interval function
            pause = false;
        });
    });
    /* **** End Banner Slider **** */


  $(".vdp-datepicker .today").trigger("click");
});



function initMap() {
    if ( $("#map").length == 0 ) return;
  
    var goodlock = { lat:40.34028, lng: -3.73678 };
  
    var map = new google.maps.Map(document.getElementById("map"), {
      center:{ lat:40.342304, lng: -3.7379202 },
      zoom: 16,
    });
  
    var contentString =
       '<div id="content" style="width:300px;color:#000;">' +
        '<div style="font-weight:500;color:#000;font-size:120%;">Zombie Bus Escape Experience</div>' +
        '<address style="color:#000;">' +
        '  CC Westfield Parquesur<br/>Av. de Gran Bretaña, S/N<br/>28916 Leganes, Madrid' +
        '  <br /> (Al lado de Mediamark)</p>' +
        '</address>' +
        '<div><b class="subway line-12"></b> Metro El Carrascal <i style="font-weight:500;">8 mins</i></div>' +
        '<div><b class="train line-c5"></b>Cercanias Zarzaquemada<i style="font-weight:500;">16 mins</i></div>' +
      "</div>";
  
    var infowindow = new google.maps.InfoWindow({
      content: contentString
    });
  
    var marker = new google.maps.Marker({
      position: goodlock,
      map: map,
      title: "Zombie Bus Escape Experience"
    });
  
    infowindow.open(map, marker);
  }

var header = document.getElementById("header");
var sticky = document.getElementById("header").offsetTop;

window.onscroll = function() {
    if (window.pageYOffset > sticky) {
      header.classList.add("sticky");
    } else {
      header.classList.remove("sticky");
    }
}; 
