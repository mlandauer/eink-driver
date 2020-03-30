#include <string.h>
#include "miniGUI.h"

extern IT8951DevInfo gstI80DevInfo;
extern uint8_t* gpFrameBuf; //Host Source Frame buffer

void EPD_Clear(uint8_t Color)
{
	memset(gpFrameBuf, Color, gstI80DevInfo.usPanelW * gstI80DevInfo.usPanelH);
}

void EPD_DrawPixel(uint16_t x0, uint16_t y0, uint8_t color)
{
	if(x0 < 0 || x0 >= gstI80DevInfo.usPanelW || y0 < 0 || y0 >= gstI80DevInfo.usPanelH)
		return ;

	/*
	�Դ洦��
	*/
	gpFrameBuf[y0*gstI80DevInfo.usPanelW + x0] = color;
}

//��ʾBMPͼƬ���ص�ȫ�ֱ���
struct   fb_var_screeninfo vinfo; 	// ��¼�û����޸ĵ���ʾ����������
uint32_t fbfd = 0;					// /dev/fb0 �ļ�������
uint8_t *fbp = 0;					// �ڴ�ӳ��ָ��
uint8_t *bmp_dst_buf = NULL;		// BMP Buffer
uint8_t *bmp_src_buf = NULL;		// BMP Buffer
uint32_t bmp_width, bmp_height;
uint8_t  bmp_BitCount;
uint32_t bytesPerLine;
uint32_t imageSize;
uint32_t skip;
PIXEL    palette[256];

static void Bitmap_format_Matrix(uint8_t *dst,uint8_t *src)
{
	uint32_t i,j,k;
    uint8_t *psrc = src;
    uint8_t *pdst = dst;
    uint8_t *p = psrc;
	uint8_t temp;
	uint32_t count;

	//����bmp�洢�ǴӺ�����ǰ�棬������Ҫ��������ת��
	switch(bmp_BitCount)
	{
		case 1:
			pdst += (bmp_width * bmp_height);

			for(i=0;i<bmp_height;i++)
			{
				pdst -= bmp_width;
				count = 0;
				for (j=0;j<(bmp_width+7)/8;j++)
				{
					temp = p[j];

					for (k=0;k<8;k++)
					{
						pdst[0]= ((temp & (0x80>>k)) >> (7-k));
						count++;
						pdst++;
						if (count == bmp_width)
						{
							break;
						}
					}
				}
				pdst -= bmp_width;
				p += bytesPerLine;
			}
		break;
		case 4:
			pdst += (bmp_width * bmp_height);

			for(i=0;i<bmp_height;i++)
			{
				pdst -= bmp_width;
				count = 0;
				for (j=0;j<(bmp_width+1)/2;j++)
				{
					temp = p[j];
					pdst[0]= ((temp & 0xf0) >> 4);
					count++;
					pdst++;
					if (count == bmp_width)
					{
						break;
					}

					pdst[0] = temp & 0x0f;
					count++;
					pdst++;
					if (count == bmp_width)
					{
						break;
					}
				}
				pdst -= bmp_width;
				p += bytesPerLine;
			}
		break;
		case 8:
			pdst += (bmp_width*bmp_height);
			for(i=0;i<bmp_height;i++)
			{
				p = psrc+(i+1)*bytesPerLine;
				p -= skip;
				for(j=0;j<bmp_width;j++)
				{
					pdst -= 1;
					p -= 1;
					pdst[0] = p[0];
				}
			}
		break;
		case 16:
			pdst += (bmp_width*bmp_height*2);
			for(i=0;i<bmp_height;i++)
			{
				p = psrc+(i+1)*bytesPerLine;
				p -= skip;
				for(j=0;j<bmp_width;j++)
				{
					pdst -= 2;
					p -= 2;
					pdst[0] = p[1];
					pdst[1] = p[0];
				}
			}
		break;
		case 24:
			pdst += (bmp_width*bmp_height*3);
			for(i=0;i<bmp_height;i++)
			{
				p = psrc+(i+1)*bytesPerLine;
				p -= skip;
				for(j=0;j<bmp_width;j++)
				{
					pdst -= 3;
					p -= 3;
					pdst[0] = p[2];
					pdst[1] = p[1];
					pdst[2] = p[0];
				}
			}
		break;
		case 32:
			pdst += (bmp_width*bmp_height*4);
			for(i=0;i<bmp_height;i++)
			{
				p = psrc+(i+1)*bmp_width*4;
				for(j=0;j<bmp_width;j++)
				{
					pdst -= 4;
					p -= 4;
					pdst[0] = p[2];
					pdst[1] = p[1];
					pdst[2] = p[0];
					pdst[3] = p[3];
				}
			}
		break;

		default:
		break;
	}
}

static void DrawMatrix(uint16_t Xpos, uint16_t Ypos,uint16_t Width, uint16_t High,const uint8_t* Matrix)
{
	uint16_t i,j,x,y;
	uint8_t R,G,B;
	uint8_t temp1,temp2;
	double Gray;

	for (y=0,j=Ypos;y<High;y++,j++)
	{
 		for (x=0,i=Xpos;x<Width;x++,i++)
		{
			switch(bmp_BitCount)
			{
				case 1:
				case 4:
				case 8:
					R = palette[Matrix[(y*Width+x)]].red;
					G = palette[Matrix[(y*Width+x)]].green;
					B = palette[Matrix[(y*Width+x)]].blue;
				break;

				case 16:
					temp1 = Matrix[(y*Width+x)*2];
					temp2 = Matrix[(y*Width+x)*2+1];
					R = (temp1 & 0x7c)<<1;
					G = (((temp1 & 0x03) << 3 ) | ((temp2&0xe0) >> 5))<<3;
					B = (temp2 & 0x1f)<<3;
				break;

				case 24:
					R = Matrix[(y*Width+x)*3];
					G = Matrix[(y*Width+x)*3+1];
					B = Matrix[(y*Width+x)*3+2];
				break;

				case 32:
					R = Matrix[(y*Width+x)*4];
					G = Matrix[(y*Width+x)*4+1];
					B = Matrix[(y*Width+x)*4+2];
				break;

				default:
				break;
			}

			Gray = (R*299 + G*587 + B*114 + 500) / 1000;
			EPD_DrawPixel(i, j, Gray);
		}
	}
}

uint8_t Show_bmp(uint32_t x, uint32_t y,char *path)
{
	FILE *fp;//bmp�ļ�ָ��
	BITMAPFILEHEADER FileHead;
	BITMAPINFOHEADER InfoHead;
	uint32_t total_length;
	uint8_t *buf = NULL;
	uint32_t ret = -1;

	fp = fopen(path,"rb");
	if (fp == NULL)
	{
		return(-1);
	}

	ret = fread(&FileHead, sizeof(BITMAPFILEHEADER),1, fp);
	if (ret != 1)
	{
		printf("Read header error!\n");
		fclose(fp);
		return(-2);
	}

	//�����Ƿ���bmpͼ��
	if (FileHead.cfType != 0x4D42)//�ļ����ͣ�"BM"(0x4D42)
	{
		printf("It's not a BMP file\n");
		fclose(fp);
		return(-3);
	}

	printf("\n*****************************************\n");
	printf("BMP_cfSize:%d \n", FileHead.cfSize);
 	printf("BMP_cfoffBits:%d \n", FileHead.cfoffBits);

	ret = fread((char *)&InfoHead, sizeof(BITMAPINFOHEADER),1, fp);
	if (ret != 1)
	{
		printf("Read infoheader error!\n");
		fclose(fp);
		return(-4);
	}

	printf("BMP_ciSize:%d \n", InfoHead.ciSize);
 	printf("BMP_ciWidth:%d \n", InfoHead.ciWidth);
	printf("BMP_ciHeight:%d \n", InfoHead.ciHeight);
	printf("BMP_ciPlanes:%x \n", InfoHead.ciPlanes);
	printf("BMP_ciBitCount:%x \n", InfoHead.ciBitCount);
	printf("BMP_ciCompress:%x \n", InfoHead.ciCompress);
	printf("BMP_ciSizeImage:%x \n", InfoHead.ciSizeImage);
	printf("BMP_ciXPelsPerMeter:%x \n", InfoHead.ciXPelsPerMeter);
	printf("BMP_ciYPelsPerMeter:%x \n", InfoHead.ciYPelsPerMeter);
	printf("BMP_ciClrUsed:%x \n", InfoHead.ciClrUsed);
	printf("BMP_ciClrImportant:%x \n", InfoHead.ciClrImportant);
 	printf("*****************************************\n\n");

	total_length = FileHead.cfSize-FileHead.cfoffBits;
	bytesPerLine=((InfoHead.ciWidth*InfoHead.ciBitCount+31)>>5)<<2;
	imageSize=bytesPerLine*InfoHead.ciHeight;
	skip=(4-((InfoHead.ciWidth*InfoHead.ciBitCount)>>3))&3;

	printf("total_length:%d,%d\n", InfoHead.ciSizeImage,FileHead.cfSize-FileHead.cfoffBits);
	printf("bytesPerLine = %d\n", bytesPerLine);
	printf("imageSize = %d\n", imageSize);
	printf("skip = %d\n", skip);

    bmp_width = InfoHead.ciWidth;
    bmp_height = InfoHead.ciHeight;
	bmp_BitCount = InfoHead.ciBitCount;

    bmp_src_buf = (uint8_t*)calloc(1,total_length);
    if(bmp_src_buf == NULL){
        printf("Load > malloc bmp out of memory!\n");
        return -5;
    }

	bmp_dst_buf = (uint8_t*)calloc(1,total_length);
    if(bmp_dst_buf == NULL){
        printf("Load > malloc bmp out of memory!\n");
        return -5;
    }

	 //��ת��������
    fseek(fp, FileHead.cfoffBits, SEEK_SET);

	//ÿ���ֽ���
    buf = bmp_src_buf;
    while ((ret = fread(buf,1,total_length,fp)) >= 0)
	{
        if (ret == 0)
		{
            usleep(100);
            continue;
        }
		buf = ((uint8_t*)buf) + ret;
        total_length = total_length - ret;
        if(total_length == 0)
            break;
    }

	//��ת�����԰�
	switch(bmp_BitCount)
	{
		case 1:
			fseek(fp, 54, SEEK_SET);
			ret = fread(palette,1,4*2,fp);
			if (ret != 8)
			{
				printf("Error: fread != 8\n");
				return -5;
			}

			bmp_dst_buf = (uint8_t*)calloc(1,InfoHead.ciWidth * InfoHead.ciHeight);
			if(bmp_dst_buf == NULL)
			{
				printf("Load > malloc bmp out of memory!\n");
				return -5;
			}
		break;

		case 4:
			fseek(fp, 54, SEEK_SET);
			ret = fread(palette,1,4*16,fp);
			if (ret != 64)
			{
				printf("Error: fread != 64\n");
				return -5;
			}

			bmp_dst_buf = (uint8_t*)calloc(1,InfoHead.ciWidth * InfoHead.ciHeight);
			if(bmp_dst_buf == NULL)
			{
				printf("Load > malloc bmp out of memory!\n");
				return -5;
			}
		break;

		case 8:
			fseek(fp, 54, SEEK_SET);

			ret = fread(palette,1,4*256,fp);

			if (ret != 1024)
			{
				printf("Error: fread != 1024\n");
				return -5;
			}
		break;

		default:
		break;
	}

	Bitmap_format_Matrix(bmp_dst_buf,bmp_src_buf);
	DrawMatrix(x, y,InfoHead.ciWidth, InfoHead.ciHeight, bmp_dst_buf);

	fclose(fp);
	return(0);
}
